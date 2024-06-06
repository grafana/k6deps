import { parse } from "k6/x/yaml";

const source = `
- name: foo
- name: bar
`;

const devices = parse(source);

function getDevice() {
  return devices[0];
}

module.exports = { getDevice };
