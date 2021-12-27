
export const samplesLeft = [
  {
    text: "unchanged input",
    code: ".",
    input_q: ".",
    input_j: '{ "foo": { "bar": { "baz": 123 } } }'
  },
  {
    text: "value at key",
    code: ".foo, .foo.bar, .foo?",
    input_q: ".foo",
    input_j: '{"foo": 42, "bar": "less interesting data"}'
  },
  {
    text: "array operation",
    code: ".[], .[]?, .[2], .[10:15]",
    input_q: ".[1]",
    input_j: '[{"name":"JSON", "good":true}, {"name":"XML", "good":false}]'
  },
  {
    text: "array/object construction",
    code: "[], {}",
    input_q: "{user, title: .titles[]}",
    input_j: '{"user":"stedolan","titles":["JQ Primer", "More JQ"]}'
  },
  {
    text: "length of a value",
    code: "length",
    input_q: ".[] | length",
    input_j: '[[1,2], "string", {"a":2}, null]'
  },
  {
    text: "keys in an array",
    code: "keys",
    input_q: "keys",
    input_j: '{"abc": 1, "abcd": 2, "Foo": 3}'
  }
];

export const samplesRight = [
  {
    text: "feed input into multiple filters",
    code: ",",
    input_q: ".foo, .bar",
    input_j: '{ "foo": 42, "bar": "something else", "baz": true}'
  },
  {
    text: "pipe output of one filter to the next filter",
    code: "|",
    input_q: ".[] | .name",
    input_j: '[{"name":"JSON", "good":true}, {"name":"XML", "good":false}]'
  },
  {
    text: "input unchanged if foo returns true",
    code: "select(foo)",
    input_q: "map(select(. >= 2))",
    input_j: '[1,5,3,0,7]'
  },
  {
    text: "invoke filter foo for each input",
    code: "map(foo)",
    input_q: "map(.+1)",
    input_j: '[1,2,3]'
  },
  {
    text: "conditionals",
    code: "if-then-else-end",
    input_q: 'if . == 0 then "zero" elif . == 1 then "one" else "many" end',
    input_j: '2'
  },
  {
    text: "string interpolation",
    code: "\\(foo)",
    input_q: '"The input was \\(.), which is one less than \\(.+1)"',
    input_j: '42'
  }
];
