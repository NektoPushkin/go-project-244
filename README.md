# Вычислитель отличий на Go

[![hexlet-check](https://github.com/NektoPushkin/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/NektoPushkin/go-project-244/actions)

[![gendiff-tests](https://github.com/NektoPushkin/go-project-244/actions/workflows/gendiff-tests.yml/badge.svg)](https://github.com/NektoPushkin/go-project-244/actions/workflows/gendiff-tests.yml)

Консольная утилита для сравнения вложенных структур (JSON, YAML)

Учебный проект Хекслета: https://ru.hexlet.io/programs/go

## Возможности

- Сравнение **вложенных** структур произвольной глубины (JSON и YAML);
- Четыре формата вывода: `stylish` (по умолчанию), `plain`, `json` (плоское JSON-дерево отличий) и `jsonpatch` (JSON Patch по [RFC 6902]);
- Порог покрытия тестами при сборке и автопроверки в GitHub Actions.

## Стек

- Go (>= 1.26)
- Библиотеки: [`urfave/cli/v3`](https://github.com/urfave/cli/v3) (CLI), [`gopkg.in/yaml.v3`](https://github.com/go-yaml/yaml) (YAML), [`fatih/color`](https://github.com/fatih/color) (цвета), [`stretchr/testify`](https://github.com/stretchr/testify) (тесты)

## Установка

Требования: Go 1.27+ 

Для `make build`/`make lint` — требуется [`golangci-lint`](https://golangci-lint.run/):

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Клонирование и сборка:

```bash
git clone https://github.com/NektoPushkin/go-project-244.git
cd go-project-244

# через Makefile (с прогоном линтера)
make build          # бинарник появится в bin/gendiff

# или напрямую без линтера
go build -o bin/gendiff ./cmd/gendiff
```

Запуск из исходников:

```bash
go run ./cmd/gendiff --help
```

## Использование

```
gendiff [опции] <config1> <config2>
```

Поддерживаемые расширения файлов: `.json`, `.yml`, `.yaml`.

Возможные опции:
   --format, -f      формат вывода: `stylish`, `plain`, `json`, `jsonpatch` (по умолчанию: `stylish`)

### Пример: stylish (по умолчанию)

```bash
./bin/gendiff testdata/fixture/json.1.json testdata/fixture/json.2.json
```

```text
{
    common: {
      + follow: false
        setting1: Value 1
      - setting2: 200
      - setting3: true
      + setting3: null
      + setting4: blah blah
      + setting5: {
            key5: value5
        }
        setting6: {
            doge: {
              - wow:
              + wow: so much
            }
            key: value
          + ops: vops
        }
    }
    group1: {
      - baz: bas
      + baz: bars
        foo: bar
      - nest: {
            key: value
        }
      + nest: str
    }
  - group2: {
        abc: 12345
        deep: {
            id: 45
        }
    }
  + group3: {
        deep: {
            id: {
                number: 45
            }
        }
        fee: 100500
    }
}
```

Если файлы идентичны, утилита выводит `no differences found`.

[Аскинема с примером использования](https://asciinema.org/a/fyiA2eoX3ZJtjBYT)

### Пример: plain

```bash
./bin/gendiff -f plain testdata/fixture/yaml.1.yaml testdata/fixture/yaml.2.yml
```

```text
Property 'common.follow' was added with value: false
Property 'common.setting2' was removed
Property 'common.setting3' was updated. From true to null
Property 'group2' was removed
Property 'group3' was added with value: [complex value]
```

[Аскинема с примером использования](https://asciinema.org/a/9E1c3xRXVRmbYD2f)

### Пример: json (дерево отличий)

```bash
./bin/gendiff --format=json testdata/fixture/json.1.json testdata/fixture/json.2.json
```

```json
{"key":"","type":"root","children":[{"key":"common","type":"nested","children":[{"key":"follow","type":"added","value2":false},{"key":"setting1","type":"unchagned","value1":"Value 1"},{"key":"setting2","type":"deleted","value1":200},{"key":"setting3","type":"changed","value1":true},{"key":"setting4","type":"added","value2":"blah blah"},{"key":"setting5","type":"added","value2":{"key5":"value5"}},{"key":"setting6","type":"nested","children":[{"key":"doge","type":"nested","children":[{"key":"wow","type":"changed","value1":"","value2":"so much"}]},{"key":"key","type":"unchagned","value1":"value"},{"key":"ops","type":"added","value2":"vops"}]}]},{"key":"group1","type":"nested","children":[{"key":"baz","type":"changed","value1":"bas","value2":"bars"},{"key":"foo","type":"unchagned","value1":"bar"},{"key":"nest","type":"changed","value1":{"key":"value"},"value2":"str"}]},{"key":"group2","type":"deleted","value1":{"abc":12345,"deep":{"id":45}}},{"key":"group3","type":"added","value2":{"deep":{"id":{"number":45}},"fee":100500}}]}
```

[Аскинема с примером использования](https://asciinema.org/a/magN9RhIC8gYx8IE)

### Пример: jsonpatch (JSON Patch, RFC 6902)

```bash
./bin/gendiff --format=jsonpatch testdata/fixture/json.1.json testdata/fixture/json.2.json
```

```json
[
    {
        "op": "add",
        "path": "/common/follow",
        "value": false
    },
    {
        "op": "remove",
        "path": "/common/setting2"
    },
    {
        "op": "replace",
        "path": "/common/setting3",
        "value": null
    },
    {
        "op": "remove",
        "path": "/group2"
    },
    {
        "op": "add",
        "path": "/group3",
        "value": {
            "fee": 100500
        }
    }
]
```

[Аскинема с примером использования](https://asciinema.org/a/T80YUoCU9hYc1DNB)

---

<details>
<summary>Автоматические тесты Хекслета</summary>

Тесты запускаются на каждый коммит. За запуск отвечает файл `.github/workflows/hexlet-check.yml` — не удаляйте и не переименовывайте ни его, ни репозиторий.

</details>

## О Хекслете

[Хекслет](https://ru.hexlet.io/) — школа программирования: авторские программы обучения с практикой, поддержкой наставников и реальными проектами, которые остаются в резюме. Этот репозиторий — один из таких проектов.
