/*
Copyright 2024 Alexander Trost All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

type Config struct {
	RGW RGW `yaml:"rgw"`
	RBD RBD `yaml:"rbd"`
}

type RGW struct {
	Realms []*Realm `yaml:"realms"`
}

type Realm struct {
	Name          string `yaml:"name"`
	Host          string `yaml:"host"`
	AccessKey     string `yaml:"accessKey"`
	SecretKey     string `yaml:"secretKey"`
	SkipTLSVerify bool   `yaml:"skipTLSVerify"`
}

type RBD struct {
	Pools []RBDPool `yaml:"pools"`
}

type RBDPool struct {
	Namespaces []string `yaml:"namespaces"`
}
