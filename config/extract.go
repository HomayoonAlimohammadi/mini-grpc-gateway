package config

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/HomayoonAlimohammadi/mini-grpc-gateway/pb/post"
)

const (
	serviceProtoFileName string = "service.proto"
)

type MethodConfig struct {
	Name    protoreflect.Name
	GoName  string
	In      *protogen.Message
	Out     *protogen.Message
	Options *post.MiniGRPCOptions
}

type ServiceConfig struct {
	Name     protoreflect.Name
	Package  string
	GoName   string
	Methods  []MethodConfig
	GPRCPort int
}

func ExtractServiceConfig(gen *protogen.Plugin, jsonConfigRaw []byte) ([]ServiceConfig, error) {
	svcConfs := []ServiceConfig{}

	jsonConf, err := readJSONConfig(jsonConfigRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON config: %w", err)
	}

	for _, file := range gen.Files {
		if file.Proto.GetName() == serviceProtoFileName {
			for _, svc := range file.Services {
				svcName := svc.Desc.Name()

				svcJsonConf, ok := jsonConf.SvcNameToConf[string(svcName)]
				if !ok {
					return nil, fmt.Errorf("failed to find service `%s` config in json file", svcName)
				}

				svcConfs = append(svcConfs, ServiceConfig{
					Name:     svcName,
					GoName:   svc.GoName,
					Methods:  extractMethodsConfig(svc),
					GPRCPort: svcJsonConf.GRPCPort,
					Package:  string(svc.Desc.ParentFile().Package()),
				})
			}
		}
	}

	return svcConfs, nil
}

func extractMethodsConfig(svc *protogen.Service) []MethodConfig {
	methConfs := make([]MethodConfig, len(svc.Methods))

	for i, m := range svc.Methods {
		methName := m.Desc.Name()

		o, ok := proto.GetExtension(m.Desc.Options(), post.E_MiniGrpcOptions).(*post.MiniGRPCOptions)
		if !ok {
			logrus.Infof("failed to cast method `%s` option to *post.MiniGRPCOption", methName)
		}

		methConfs[i] = MethodConfig{
			Name:    methName,
			GoName:  m.GoName,
			In:      m.Input,
			Out:     m.Output,
			Options: o,
		}
	}

	return methConfs
}

type jsonConfig struct {
	SvcNameToConf map[string]*serviceJSONConfig `json:"-"`
}

type serviceJSONConfig struct {
	GRPCPort int `json:"grpc-port"`
}

func readJSONConfig(jsonConfigRaw []byte) (*jsonConfig, error) {
	c := jsonConfig{}
	err := json.Unmarshal(jsonConfigRaw, &c.SvcNameToConf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &c, nil
}
