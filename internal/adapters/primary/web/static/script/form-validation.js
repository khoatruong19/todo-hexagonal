const formValidationListener = (inputIds) => {
  document.addEventListener("htmx:afterRequest", function (event) {
    let validInputIds = inputIds;
    if (event.detail.xhr.status === 400) {
      console.log({ validInputIds, text: event.detail });
      const errors = JSON.parse(event.detail.xhr.responseText);
      errors.forEach((error) => {
        const inputId =
          error.field.charAt(0).toLowerCase() + error.field.slice(1);
        validInputIds = validInputIds.filter((id) => id !== inputId);
        const field = document.getElementById(inputId);
        if (field) field.classList.add("input-error");

        const fieldErrorMessage = document.getElementById(inputId + "Error");
        if (fieldErrorMessage) {
          fieldErrorMessage.classList.remove("hidden");
          fieldErrorMessage.innerText = error.message;
        }
      });
    }

    validInputIds.forEach((id) => {
      const field = document.getElementById(id);
      if (field && field.classList.contains("input-error"))
        field.classList.remove("input-error");

      const fieldErrorMessage = document.getElementById(id + "Error");
      if (fieldErrorMessage) fieldErrorMessage.classList.add("hidden");
    });
  });
};
