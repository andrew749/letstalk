export function onChange(model, event) {
  let fieldName = event.target.name;
  let fieldValue = event.target.value;
  console.log(fieldName + "=>" + fieldValue);
  this.setState(
    prevState => ({
      [model]: {
        ...prevState[model],
        [fieldName]: fieldValue
      }
    })
  );
}
