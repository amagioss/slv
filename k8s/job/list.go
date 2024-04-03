package main

import (
	"context"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"oss.amagi.com/slv/core/config"
	slvv1 "oss.amagi.com/slv/k8s/api/v1"
)

func listSLVs(cfg *rest.Config) ([]slvv1.SLV, error) {
	dynamicClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}
	unstructuredList, err := dynamicClient.Resource(
		schema.GroupVersionResource{
			Group:    config.K8SLVGroup,
			Version:  config.K8SLVVersion,
			Resource: "slvs",
		}).Namespace(getNamespace()).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(unstructuredList.UnstructuredContent())
	if err != nil {
		return nil, err
	}
	var slvObjList slvv1.SLVList
	if err = json.Unmarshal(jsonBytes, &slvObjList); err != nil {
		return nil, err
	}
	return slvObjList.Items, nil
}
