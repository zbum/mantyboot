package configuration

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/zbum/mantyboot/errors"
)

type Validator interface {
	Validate() error
}

type ValidationRule struct {
	Field     string
	Required  bool
	Min       *int
	Max       *int
	MinLength *int
	MaxLength *int
	Pattern   *regexp.Regexp
	Custom    func(interface{}) error
}

type ConfigurationValidator struct {
	rules map[string]ValidationRule
}

func NewConfigurationValidator() *ConfigurationValidator {
	return &ConfigurationValidator{
		rules: make(map[string]ValidationRule),
	}
}

func (cv *ConfigurationValidator) AddRule(field string, rule ValidationRule) {
	cv.rules[field] = rule
}

func (cv *ConfigurationValidator) Validate(config interface{}) error {
	val := reflect.ValueOf(config)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.WrapConfigurationError(nil, "configuration must be a struct")
	}

	var validationErrors []error

	for fieldName, rule := range cv.rules {
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			continue
		}

		if err := cv.validateField(field, rule); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	if len(validationErrors) > 0 {
		return &errors.ValidationError{
			Field:   "configuration",
			Message: fmt.Sprintf("validation failed: %v", validationErrors),
		}
	}

	return nil
}

func (cv *ConfigurationValidator) validateField(field reflect.Value, rule ValidationRule) error {
	// Check if field is zero value when required
	if rule.Required && field.IsZero() {
		return &errors.ValidationError{
			Field:   rule.Field,
			Message: "field is required",
		}
	}

	// Skip validation if field is zero and not required
	if field.IsZero() {
		return nil
	}

	// Custom validation
	if rule.Custom != nil {
		if err := rule.Custom(field.Interface()); err != nil {
			return &errors.ValidationError{
				Field:   rule.Field,
				Message: err.Error(),
			}
		}
	}

	// Type-specific validation
	switch field.Kind() {
	case reflect.String:
		return cv.validateString(field.String(), rule)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cv.validateInt(field.Int(), rule)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return cv.validateUint(field.Uint(), rule)
	case reflect.Float32, reflect.Float64:
		return cv.validateFloat(field.Float(), rule)
	}

	return nil
}

func (cv *ConfigurationValidator) validateString(value string, rule ValidationRule) error {
	if rule.MinLength != nil && len(value) < *rule.MinLength {
		return &errors.ValidationError{
			Field:   rule.Field,
			Message: fmt.Sprintf("length must be at least %d", *rule.MinLength),
		}
	}

	if rule.MaxLength != nil && len(value) > *rule.MaxLength {
		return &errors.ValidationError{
			Field:   rule.Field,
			Message: fmt.Sprintf("length must be at most %d", *rule.MaxLength),
		}
	}

	if rule.Pattern != nil && !rule.Pattern.MatchString(value) {
		return &errors.ValidationError{
			Field:   rule.Field,
			Message: fmt.Sprintf("must match pattern %s", rule.Pattern.String()),
		}
	}

	return nil
}

func (cv *ConfigurationValidator) validateInt(value int64, rule ValidationRule) error {
	if rule.Min != nil && value < int64(*rule.Min) {
		return &errors.ValidationError{
			Field:   rule.Field,
			Message: fmt.Sprintf("must be at least %d", *rule.Min),
		}
	}

	if rule.Max != nil && value > int64(*rule.Max) {
		return &errors.ValidationError{
			Field:   rule.Field,
			Message: fmt.Sprintf("must be at most %d", *rule.Max),
		}
	}

	return nil
}

func (cv *ConfigurationValidator) validateUint(value uint64, rule ValidationRule) error {
	if rule.Min != nil && value < uint64(*rule.Min) {
		return &errors.ValidationError{
			Field:   rule.Field,
			Message: fmt.Sprintf("must be at least %d", *rule.Min),
		}
	}

	if rule.Max != nil && value > uint64(*rule.Max) {
		return &errors.ValidationError{
			Field:   rule.Field,
			Message: fmt.Sprintf("must be at most %d", *rule.Max),
		}
	}

	return nil
}

func (cv *ConfigurationValidator) validateFloat(value float64, rule ValidationRule) error {
	if rule.Min != nil && value < float64(*rule.Min) {
		return &errors.ValidationError{
			Field:   rule.Field,
			Message: fmt.Sprintf("must be at least %d", *rule.Min),
		}
	}

	if rule.Max != nil && value > float64(*rule.Max) {
		return &errors.ValidationError{
			Field:   rule.Field,
			Message: fmt.Sprintf("must be at most %d", *rule.Max),
		}
	}

	return nil
}

// Built-in validation functions
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}

func ValidateHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}

	// Basic hostname validation
	hostnameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
	if !hostnameRegex.MatchString(hostname) {
		return fmt.Errorf("invalid hostname format")
	}

	return nil
}

func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Basic URL validation
	urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	if !urlRegex.MatchString(url) {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

// Tag-based validation
func ValidateStruct(config interface{}) error {
	val := reflect.ValueOf(config)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.WrapConfigurationError(nil, "configuration must be a struct")
	}

	var validationErrors []error
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if err := validateFieldByTags(field, fieldType); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	if len(validationErrors) > 0 {
		return &errors.ValidationError{
			Field:   "configuration",
			Message: fmt.Sprintf("validation failed: %v", validationErrors),
		}
	}

	return nil
}

func validateFieldByTags(field reflect.Value, fieldType reflect.StructField) error {
	tag := fieldType.Tag.Get("validate")
	if tag == "" {
		return nil
	}

	rules := strings.Split(tag, ",")
	fieldName := fieldType.Name

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)

		switch {
		case rule == "required":
			if field.IsZero() {
				return &errors.ValidationError{
					Field:   fieldName,
					Message: "field is required",
				}
			}
		case strings.HasPrefix(rule, "min="):
			if err := validateMin(field, rule, fieldName); err != nil {
				return err
			}
		case strings.HasPrefix(rule, "max="):
			if err := validateMax(field, rule, fieldName); err != nil {
				return err
			}
		case strings.HasPrefix(rule, "pattern="):
			if err := validatePattern(field, rule, fieldName); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateMin(field reflect.Value, rule, fieldName string) error {
	minStr := strings.TrimPrefix(rule, "min=")
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return errors.WrapConfigurationError(err, "invalid min value in validation tag")
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) < min {
			return &errors.ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("length must be at least %d", min),
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() < int64(min) {
			return &errors.ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("must be at least %d", min),
			}
		}
	}

	return nil
}

func validateMax(field reflect.Value, rule, fieldName string) error {
	maxStr := strings.TrimPrefix(rule, "max=")
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return errors.WrapConfigurationError(err, "invalid max value in validation tag")
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) > max {
			return &errors.ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("length must be at most %d", max),
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() > int64(max) {
			return &errors.ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("must be at most %d", max),
			}
		}
	}

	return nil
}

func validatePattern(field reflect.Value, rule, fieldName string) error {
	patternStr := strings.TrimPrefix(rule, "pattern=")
	pattern := regexp.MustCompile(patternStr)

	if field.Kind() == reflect.String && !pattern.MatchString(field.String()) {
		return &errors.ValidationError{
			Field:   fieldName,
			Message: fmt.Sprintf("must match pattern %s", patternStr),
		}
	}

	return nil
}
