package alice

import "encoding/json"

// Intent структура прототипа интента в запросе.
type Intent struct {
	Slots map[string]Entity `json:"slots"`
}

// Intents возвращает все необработанные интенты из запроса.
func (req *Request) Intents() []NamedIntent {
	res := make([]NamedIntent, 0, len(req.Request.NLU.Intents))
	for name, in := range req.Request.NLU.Intents {
		res = append(res, NamedIntent{name: name, intent: in})
	}
	return res
}

// NamedIntent интент с его именем (ключом в карте интентов).
type NamedIntent struct {
	name   string
	intent Intent
}

func (ni NamedIntent) Name() string {
	return ni.name
}

// IntentSlot слот интента с именем и сущностью.
type IntentSlot struct {
	name   string
	entity Entity
}

// Name возвращает имя слота.
func (s IntentSlot) Name() string { return s.name }

// Type возвращает тип сущности в слоте.
func (s IntentSlot) Type() string { return s.entity.Type }

// Tokens возвращает границы токенов в исходной реплике.
func (s IntentSlot) Tokens() (int, int) { return s.entity.Tokens.Start, s.entity.Tokens.End }

// Entity возвращает обработанную сущность слота.
func (s IntentSlot) Entity() (interface{}, bool) {
	if s.entity.Value == nil {
		return nil, false
	}
	h := holder(s.entity.Type)
	if h == nil {
		return nil, false
	}
	if err := json.Unmarshal(*s.entity.Value, h); err != nil {
		return nil, false
	}
	return h, true
}

// Slots возвращает все слоты данного интента.
func (ni NamedIntent) Slots() []IntentSlot {
	res := make([]IntentSlot, 0, len(ni.intent.Slots))
	for name, e := range ni.intent.Slots {
		res = append(res, IntentSlot{name: name, entity: e})
	}
	return res
}
