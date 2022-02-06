
// Basics (from tutorial)
// https://jmespath.org/tutorial.html
export const samplesJmesLeft = [
  {
    text: "basic select",
    code: `a`,
    input_q: `a`,
    input_j: `{"a": "foo", "b": "bar", "c": "baz"}`
  },
  {
    text: "nested select",
    code: `a.b.c.d`,
    input_q: `a.b.c.d`,
    input_j: `{"a": {"b": {"c": {"d": "value"}}}}`
  },
  {
    text: "index expression",
    code: `[1]`,
    input_q: `[1]`,
    input_j: `["a", "b", "c", "d", "e", "f"]`
  },
  {
    text: "slice",
    code: `[0:5]`,
    input_q: `[0:5]`,
    input_j: `[0, 1, 2, 3, 4, 5, 6, 7, 8, 9]`
  },
  {
    text: "list projection",
    code: `people[*].first`,
    input_q: `people[*].first`,
    input_j: `{
      "people": [
        {"first": "James", "last": "d"},
        {"first": "Jacob", "last": "e"},
        {"first": "Jayden", "last": "f"},
        {"missing": "different"}
      ],
      "foo": {"bar": "baz"}
    }`
  },
  {
    text: "slice projection",
    code: `people[:2].first`,
    input_q: `people[:2].first`,
    input_j: `{
      "people": [
        {"first": "James", "last": "d"},
        {"first": "Jacob", "last": "e"},
        {"first": "Jayden", "last": "f"},
        {"missing": "different"}
      ],
      "foo": {"bar": "baz"}
    }`
  },

];


// Bit more advanced (from examples page)
// https://jmespath.org/examples.html
export const samplesJmesRight = [
  {
    text: "object projection",
    code: `ops.*.numArgs`,
    input_q: `ops.*.numArgs`,
    input_j: `{
      "ops": {
        "functionA": {"numArgs": 2},
        "functionB": {"numArgs": 3},
        "functionC": {"variadic": true}
      }
    }`
  },
  {
    text: "flatten projection",
    code: `reservations[*].instances[*].state`,
    input_q: `reservations[*].instances[*].state`,
    input_j: `{
      "reservations": [
        {
          "instances": [
            {"state": "running"},
            {"state": "stopped"}
          ]
        },
        {
          "instances": [
            {"state": "terminated"},
            {"state": "running"}
          ]
        }
      ]
    }`
  },
  {
    text: "filter projections",
    code: `machines[?state=='running'].name`,
    input_q: `machines[?state=='running'].name`,
    input_j: `{
      "machines": [
        {"name": "a", "state": "running"},
        {"name": "b", "state": "stopped"},
        {"name": "b", "state": "running"}
      ]
    }`
  },
  {
    text: "pipe",
    code: `people[*].first | [0]`,
    input_q: `people[*].first | [0]`,
    input_j: `{
      "people": [
        {"first": "James", "last": "d"},
        {"first": "Jacob", "last": "e"},
        {"first": "Jayden", "last": "f"},
        {"missing": "different"}
      ],
      "foo": {"bar": "baz"}
    }`
  },
  {
    text: "multiselect array",
    code: `people[].[name, state.name]`,
    input_q: `people[].[name, state.name]`,
    input_j: `{
      "people": [
        {
          "name": "a",
          "state": {"name": "up"}
        },
        {
          "name": "b",
          "state": {"name": "down"}
        },
        {
          "name": "c",
          "state": {"name": "up"}
        }
      ]
    }`
  },
  {
    text: "multiselect object",
    code: `people[].{Name: name, State: state.name}`,
    input_q: `people[].{Name: name, State: state.name}`,
    input_j: `{
      "people": [
        {
          "name": "a",
          "state": {"name": "up"}
        },
        {
          "name": "b",
          "state": {"name": "down"}
        },
        {
          "name": "c",
          "state": {"name": "up"}
        }
      ]
    }`
  },

];
