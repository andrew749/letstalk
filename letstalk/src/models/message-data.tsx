export default class MessageData {

  _user: string;
  _id: number;
  _onClickAction: () => any;

  constructor(id: number, user: string) {
    this._user = user;
    this._id = id;
  }

  getUserName(): string {
    return this._user;
  }

  getId(): number {
    return this._id;
  }

  setOnClickAction(action: () => any) {
    this._onClickAction = action;
  }

  onPress() {
    this._onClickAction();
  }

}
