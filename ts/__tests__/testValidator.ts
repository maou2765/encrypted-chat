/**
 * @jest-environment jsdom
 */
import validator from "../src/validator";
import { screen } from "@testing-library/dom";
import "@testing-library/jest-dom";

const HTMLForm = `<html>
<form>
<input type="text" name="given_name" value="abc a" data-testid="given_name"/>
<input type="text" name="surnname" value="abc" data-testid="surnname"/>/>
<input type="text" name="middle_name" value="abc" data-testid="middle_name"/>/>
<input type="text" name="long_name" value="abcdefghi" data-testid="long_name"/>/>
<input type="text" name="short_name" value="ab" data-testid="short_name"/>/>
<input type="text" name="no_name" value="" data-testid="no_name"/>/>
<input type="text" name="valid_email" value="abc@email.com" data-testid="valid_email"/>/>
<input type="text" name="invalid_email" value="abcemail.com" data-testid="invalid_email"/>/>
<input type="number" name="large_number" value="12" data-testid="large_number"/>/>
<input type="number" name="small_number" value="1" data-testid="small_number"/>/>
</form>
</html>`;
test.each([
  [{ given_name: "required" }, true],
  [{ no_name: "required" }, false],
  [{ given_name: "required|alpha" }, false],
  [{ given_name: "alpha" }, false],
  [{ long_name: "alpha" }, true],
  [{ long_name: "required|alpha" }, true],
  [{ given_name: "maxLength:5" }, true],
  [{ short_name: "minLength:5" }, false],
  [{ long_name: "minLength:5" }, true],
  [{ long_name: "regexp:^[a-z]*$" }, true],
  [{ long_name: "regexp:^[abc]*$" }, false],
  [{ given_name: "required|maxLength:5" }, true],
  [{ given_name: "required|maxLength:4" }, false],
  [{ given_name: "required|maxLength:5|same:surnname" }, false],
  [{ surnname: "required|maxLength:5|same:middle_name" }, true],
  [{ large_number: "max:5" }, false],
  [{ large_number: "max:12" }, true],
  [{ small_number: "min:1" }, true],
  [{ small_number: "min:2" }, false],
  [{ valid_email: "email" }, true],
  [{ invalid_email: "email" }, false],
  [{ valid_email: "email|alpha", large_number: "max:12" }, false],
  [
    { valid_email: "email|required|maxLength:20", large_number: "max:12" },
    true,
  ],
])("condition %s is %s", (cond, res) => {
  document.body.innerHTML = HTMLForm;
  expect(validator(cond)).toBe(res);
});
test.each([
  [{ given_name: "required" }, true],
  [{ no_name: "required" }, false],
  [{ given_name: "required|alpha" }, false],
  [{ given_name: "alpha" }, false],
  [{ long_name: "alpha" }, true],
  [{ long_name: "required|alpha" }, true],
  [{ given_name: "maxLength:5" }, true],
  [{ short_name: "minLength:5" }, false],
  [{ long_name: "minLength:5" }, true],
  [{ long_name: "regexp:^[a-z]*$" }, true],
  [{ long_name: "regexp:^[abc]*$" }, false],
  [{ given_name: "required|maxLength:5" }, true],
  [{ given_name: "required|maxLength:4" }, false],
  [{ given_name: "required|maxLength:5|same:surnname" }, false],
  [{ surnname: "required|maxLength:5|same:middle_name" }, true],
  [{ large_number: "max:5" }, false],
  [{ large_number: "max:12" }, true],
  [{ small_number: "min:1" }, true],
  [{ small_number: "min:2" }, false],
  [{ valid_email: "email" }, true],
  [{ invalid_email: "email" }, false],
  [{ valid_email: "email|alpha", large_number: "max:12" }, false],
  [
    { valid_email: "email|required|maxLength:20", large_number: "max:12" },
    true,
  ],
  ,
])("condition %s outlook test", (cond, res) => {
  document.body.innerHTML = HTMLForm;
  test.each(
    Object.keys(cond).map((key) => {
      if (
        validator({
          [key]: cond[key],
        }) === false
      ) {
        expect(screen.getByTestId(key)).toHaveClass("is-invalid");
      } else {
        expect(screen.getByTestId(key)).not.toHaveClass("is-invalid");
      }
    })
  );
});
