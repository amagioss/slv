package providers

// type ProviderInput struct {
// 	inputs map[string]any
// }

// func (pi *ProviderInput) GetInput(name string) any {
// 	return pi.inputs[name]
// }

// func (pi *ProviderInput) NewSecretKey() (secretKey *crypto.SecretKey, err error) {
// 	return crypto.NewSecretKey(environments.EnvironmentKey)
// }

// func (pi *ProviderInput) NewSecretKeyFromPassword(password []byte) (secretKey *crypto.SecretKey, salt []byte, err error) {
// 	return crypto.NewSecretKeyFromPassword(password, environments.EnvironmentKey)
// }

// func (pi *ProviderInput) NewSecretKeyFromPasswordWithoutSalt(password []byte) (secretKey *crypto.SecretKey, err error) {
// 	return crypto.NewSecretKeyFromPasswordAndSalt(password, nil, environments.EnvironmentKey)
// }

// type ProviderInputRequirements struct {
// 	inputRequirements map[string]inputMeta
// }

// type InputType string

// type inputMeta struct {
// 	description string
// 	required    bool
// 	inputType   InputType
// }

// func (pir *ProviderInputRequirements) AddInputRequirement(name, description string, required bool, inputType InputType) {
// 	if pir.inputRequirements == nil {
// 		pir.inputRequirements = make(map[string]inputMeta)
// 	}
// 	pir.inputRequirements[name] = inputMeta{
// 		description: description,
// 		required:    required,
// 		inputType:   inputType,
// 	}
// }

// func (pir *ProviderInputRequirements) GetInputRequirements(name string) (description string, required bool, inputType *InputType) {
// 	if pir.inputRequirements != nil {
// 		if input, ok := pir.inputRequirements[name]; ok {
// 			return input.description, input.required, &input.inputType
// 		}
// 	}
// 	return "", false, nil
// }

// func (pir *ProviderInputRequirements) GetInputNames() []string {
// 	if pir.inputRequirements != nil {
// 		inputs := make([]string, 0, len(pir.inputRequirements))
// 		for input := range pir.inputRequirements {
// 			inputs = append(inputs, input)
// 		}
// 		return inputs
// 	}
// 	return nil
// }
