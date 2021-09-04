import validator from "./utils/validator";
import * as $ from "jquery";

$(document).ready(function () {
  $("#submit-btn").click(function (e) {
    e.preventDefault();
    e.stopPropagation();
    const result = validator({
      given_name: "required|maxLength:255",
      surname: "maxLength:255",
      bio: "maxLength:255",
      email: "required|email",
      password: "required",
      confirm_password: "same:password",
    });
    if (result) {
      $("form[action|='/signup']").trigger("submit");
    }
  });
});
