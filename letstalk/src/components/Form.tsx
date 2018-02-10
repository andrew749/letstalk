import { InjectedFormProps } from 'redux-form';

interface FormOnSubmitProp<FormData> {
  onSubmit(values: FormData): void;
}

// Additional props added to form, passed into the P type parameter for redux-form types.
export type FormP<FormData, P = {}> = FormOnSubmitProp<FormData> & P

export type FormProps<FormData, P = {}> =
  FormP<FormData, P> & InjectedFormProps<FormData, FormP<FormData, P>>;
