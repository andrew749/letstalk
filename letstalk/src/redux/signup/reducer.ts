import { ActionTypes, TypeKeys } from './actions';

export type State = {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  phoneNumber: string;
  errorMsg: string | null;
};

const initialState: State = {
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  phoneNumber: '',
  errorMsg: null,
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.RECEIVED_ERROR_MSG:
      return {
        ...state,
        errorMsg: action.value,
      };
    case TypeKeys.SET_FIRST_NAME:
      return {
        ...state,
        firstName: action.value,
      };
    case TypeKeys.SET_LAST_NAME:
      return {
        ...state,
        lastName: action.value,
      };
    case TypeKeys.SET_EMAIL:
      return {
        ...state,
        email: action.value,
      };
    case TypeKeys.SET_PASSWORD:
      return {
        ...state,
        password: action.value,
      };
    case TypeKeys.SET_PHONE_NUMBER:
      return {
        ...state,
        phoneNumber: action.value,
      };
    case TypeKeys.RESET:
      return initialState;
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
      return state;
  }
};
