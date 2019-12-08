package elasticsearch

// IndexSettings is the index settings and mappings.
const IndexSettings string = `
{
	"settings": {
	  "analysis": {
		"char_filter": {
		  "zero_width_spaces": {
			  "type":       "mapping",
			  "mappings": [ "\\u200C=>\\u0020"] 
		  }
		},
		"filter": {
		  "autocomplete": {
			"type": "edge_ngram",
			"min_gram": 1,
			"max_gram": 20
		  },
		  "arabic_stop": {
			"type":       "stop",
			"stopwords":  "_arabic_" 
		  },
		  "arabic_stemmer": {
			"type":       "stemmer",
			"language":   "arabic"
		  },
		  "persian_stop": {
			"type":       "stop",
			"stopwords":  "_persian_" 
		  },
		  "english_stop": {
			"type":       "stop",
			"stopwords":  "_english_" 
		  },
		  "english_stemmer": {
			"type":       "stemmer",
			"language":   "english"
		  },
		  "english_possessive_stemmer": {
			"type":       "stemmer",
			"language":   "possessive_english"
		  },
		  "words_splitter": {
			"type": "word_delimiter",
			"preserve_original": "true"
		  }
		},
		"analyzer": {
		  "nameAnalyzer": {
			"tokenizer": "standard",
			"filter": [
			  "lowercase",
			  "autocomplete",
			  "words_splitter"
			]
		  },
		  "rebuilt_english": {
			"tokenizer":  "standard",
			"filter": [
			  "english_possessive_stemmer",
			  "lowercase",
			  "english_stop",
			  "english_stemmer",
			  "autocomplete"
			]
		  },
		  "rebuilt_arabic": {
			"tokenizer":  "standard",
			"filter": [
			  "lowercase",
			  "decimal_digit",
			  "arabic_stop",
			  "arabic_normalization",
			  "arabic_stemmer",
			  "autocomplete"
			]
		  },
		  "rebuilt_persian": {
			"tokenizer":     "standard",
			"char_filter": [ "zero_width_spaces" ],
			"filter": [
			  "lowercase",
			  "decimal_digit",
			  "arabic_normalization",
			  "persian_normalization",
			  "persian_stop",
			  "autocomplete"
			]
		  }
		}
	  }
	},
	"mappings": { 
	  "properties": {
		"createdAt": {
		  "type": "long"
		},
		"id": {
		  "type": "text",
		  "fields": {
			"keyword": {
			  "type": "keyword",
			  "ignore_above": 256
			}
		  }
		},
		"name": {
		  "type": "text",
		  "fields": {
			"keyword": {
			  "type": "keyword",
			  "ignore_above": 256
			},
			"text": {
			  "type": "text",
			  "analyzer": "nameAnalyzer",
			  "search_analyzer": "standard"
			}
		  }
		},
		"ownerID": {
		  "type": "text",
		  "fields": {
			"keyword": {
			  "type": "keyword",
			  "ignore_above": 256
			}
		  }
		},
		"size": {
		  "type": "long"
		},
		"type": {
		  "type": "text",
		  "fields": {
			"keyword": {
			  "type": "keyword",
			  "ignore_above": 256
			}
		  }
		},
		"updatedAt": {
		  "type": "long"
		}
	  }
	}
  }
`
