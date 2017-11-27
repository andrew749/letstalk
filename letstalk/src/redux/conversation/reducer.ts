import MessageData from '../../models/message-data';

// TODO: Add a type
const initialState = {
  threads: [
    new MessageData(1, "Andrew"),
    new MessageData(2, "Adam"),
  ],
};

export default function conversationsReducer(state = initialState, action: any) {
  return state;
}