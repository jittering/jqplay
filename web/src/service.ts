import axios from "axios";

export default class Service {

  getJqVersion() {
    return axios.get("/jq/version").then((res) => res.data);
  }

}
