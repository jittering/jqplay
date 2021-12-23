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
}

export default class Service {
  getJqInput() {
    return axios.get("/jq/input").then((res) => res.data);
  }

  getJqVersion() {
    return axios.get("/jq/version").then((res) => res.data);
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

}
