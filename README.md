# Mox
a extension for github.com/samber/mo.

## Validate
- [x] github.com/go-playground/validator 
  - RegisterGPValidatorNotNil: notnil: mandatory, allows zero value (except nil)
  - RegisterGPValidatorPresent: present, require option.IsPresent=true
  - RegisterGPVUnwrapOptionTypeFunc: to unwrap option value, make it value pass to next validate tag.

## Web
- github.com/gin-gonic/gin
  - [x] form: OptionFormBinding
  - [x] query: OptionQueryBinding
  - [x] uri: ShouldBindGinUri


