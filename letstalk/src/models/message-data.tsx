export default class MessageData {

  _user: string;
  _id: number;
  _onClickAction: () => any;

  constructor(id, user) {
    this._user = user;
    this._id = id;
  }

  getUserName(): string {
    return this._user;
  }

  getId(): number {
    return this._id;
  }

  setOnClickAction(action) {
    this._onClickAction = action;
  }

  onPress() {
    this._onClickAction();
  }

}
