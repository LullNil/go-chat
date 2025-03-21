const keytar = require("keytar");

const SERVICE_NAME = "GoChat";
const ACCOUNT_NAME = "userToken";

module.exports = {
  async saveToken(token) {
    await keytar.setPassword(SERVICE_NAME, ACCOUNT_NAME, token);
  },
  async getToken() {
    return await keytar.getPassword(SERVICE_NAME, ACCOUNT_NAME);
  },
  async deleteToken() {
    await keytar.deletePassword(SERVICE_NAME, ACCOUNT_NAME);
  },
};
