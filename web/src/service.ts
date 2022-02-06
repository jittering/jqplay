import axios from "axios";

export interface JqOpt {
  name: string
  enabled: boolean
}

export class Jq {
  j: string // input json
  q: string // filter query
  o: Array<JqOpt> // options

  constructor() {
    this.j = "";
    this.q = ".";
    this.o = [
      { name: "slurp", enabled: false },
      { name: "null-input", enabled: false },
      { name: "compact-output", enabled: false },
      { name: "raw-input", enabled: false },
      { name: "raw-output", enabled: false },
    ];
  }

  getOpt(name) {
    return this.o.find((el) => el.name === name);
  }
}

export default class Service {
  getJqInput() {
    return axios.get("/jq/input").then((res) => res.data);
  }

  getJqVersion() {
    return axios.get("/jq/version").then((res) => res.data);
  }

  getJqCommandLine(data: Jq) {
    return axios.post("/jq/commandline", {
      j: data.j,
      q: data.q,
      o: data.o
    }).then((res) => res.data);
  }

  runJq(data: Jq) {
    return axios.post("/jq", {
      j: data.j,
      q: data.q,
      o: data.o
    },
      {
        // return the val directly, don't do any json parsing on the result
        transformResponse: (value) => value
      },
    ).then((res) => res.data);
  }

  runJmesPath(data: Jq): Promise<string | void> {
    return axios.post("/jmespath", {
      input: data.j,
      search: data.q,
    },
      {
        transformResponse: (value) => value
      },
    ).then((res) => res.data);
  }

}
