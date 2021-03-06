package agents

import (
	"fmt"
	"github.com/TeaWeb/code/teaconfigs"
	"github.com/iwind/TeaGo/types"
	"github.com/iwind/TeaGo/utils/string"
	"regexp"
	"strings"
)

// 阈值定义
type Threshold struct {
	Id            string            `yaml:"id" json:"id"`                       // ID
	Param         string            `yaml:"param" json:"param"`                 // 参数
	Operator      ThresholdOperator `yaml:"operator" json:"operator"`           // 运算符
	Value         string            `yaml:"value" json:"value"`                 // 对比值
	NoticeLevel   NoticeLevel       `yaml:"noticeLevel" json:"noticeLevel"`     // 通知级别
	NoticeMessage string            `yaml:"noticeMessage" json:"noticeMessage"` // 通知消息

	regValue   *regexp.Regexp
	floatValue float64
}

// 新阈值对象
func NewThreshold() *Threshold {
	return &Threshold{
		Id: stringutil.Rand(16),
	}
}

// 校验
func (this *Threshold) Validate() error {
	if this.Operator == ThresholdOperatorRegexp {
		reg, err := regexp.Compile(this.Value)
		if err != nil {
			return err
		}
		this.regValue = reg
	} else if this.Operator == ThresholdOperatorGt || this.Operator == ThresholdOperatorGte || this.Operator == ThresholdOperatorLt || this.Operator == ThresholdOperatorLte {
		this.floatValue = types.Float64(this.Value)
	}
	return nil
}

// 将此条件应用于阈值，检查是否匹配
func (this *Threshold) Test(value interface{}) bool {
	paramValue := teaconfigs.RegexpNamedVariable.ReplaceAllStringFunc(this.Param, func(s string) string {
		if value == nil {
			return ""
		}

		varName := s[2 : len(s)-1]
		switch v := value.(type) {
		case string:
			if varName == "0" {
				return v
			}
			return ""
		case int8, int16, int, int32, int64, uint8, uint16, uint, uint32, uint64:
			if varName == "0" {
				return fmt.Sprintf("%d", v)
			}
			return "0"
		case float32, float64:
			if varName == "0" {
				return fmt.Sprintf("%f", v)
			}
			return "0"
		case bool:
			if varName == "0" {
				if v {
					return "1"
				}
				return "0"
			}
			return "0"
		case []interface{}:
			index := types.Int(varName)
			if index >= 0 && index < len(v) {
				return types.String(v[index])
			}
			return ""
		case map[string]interface{}:
			result, found := v[varName]
			if found {
				return types.String(result)
			}
			return ""
		}
		return s
	})

	switch this.Operator {
	case ThresholdOperatorRegexp:
		if this.regValue == nil {
			return false
		}
		return this.regValue.MatchString(types.String(paramValue))
	case ThresholdOperatorGt:
		return types.Float64(paramValue) > this.floatValue
	case ThresholdOperatorGte:
		return types.Float64(paramValue) >= this.floatValue
	case ThresholdOperatorLt:
		return types.Float64(paramValue) < this.floatValue
	case ThresholdOperatorLte:
		return types.Float64(paramValue) <= this.floatValue
	case ThresholdOperatorEq:
		return paramValue == this.Value
	case ThresholdOperatorNot:
		return paramValue != this.Value
	case ThresholdOperatorPrefix:
		return strings.HasPrefix(types.String(paramValue), this.Value)
	case ThresholdOperatorSuffix:
		return strings.HasSuffix(types.String(paramValue), this.Value)
	case ThresholdOperatorContains:
		return strings.Contains(types.String(paramValue), this.Value)
	}
	return false
}
