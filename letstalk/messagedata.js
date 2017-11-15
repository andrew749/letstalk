export default class MessageData {

  constructor(id, user) {
    this._user = user;
    this._id = id;
  }

  getUserName() {
    return this._user;
  }

  getId() {
    return this._id;
  }

}
