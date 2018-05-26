import { Dispatch } from 'react-redux';
import { ToastActionsCreators } from 'react-native-redux-toast';

import { RootState } from '../redux';

export function errorToast(message: string) {
  return (dispatch: Dispatch<RootState>): Promise<void> => {
    return dispatch(ToastActionsCreators.displayError(message, 2000));
  };
}

export function infoToast(message: string) {
  return (dispatch: Dispatch<RootState>): Promise<void> => {
    return dispatch(ToastActionsCreators.displayInfo(message, 2000));
  };
}
