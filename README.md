# Mox
a extension for github.com/samber/mo.

## Validate
- [x] github.com/go-playground/validator 
  - RegisterGPValidatorNotNil: add json tag notnil, mandatory, allows zero value (except nil)
  - RegisterGPValidatorPresent: add json tag present, require option.IsPresent=true
  - RegisterGPVUnwrapOptionTypeFunc: to unwrap option value, make it value pass to next validate tag.

## Web
- github.com/gin-gonic/gin
  - [x] form: use OptionFormBinding
  - [x] query: use OptionQueryBinding
  - [x] uri: call ShouldBindGinUri
  - [ ] add json field for form、query.

## Json
reason: for https://github.com/samber/mo/pull/65 trust set null as set a value.  

solution: 
1. upgrade to go1.24 to use omitzero.
2. use other json library to ignore serialize the field which IsPresent=false.
   - here use github.com/json-iterator/go
   - Usage: register  OptionExtension, and add json tag omitempty.



