import * as $ from "jquery";
import isEmail from "validator/lib/isEmail";
/**
 * This function is used with jQuery to get form value.
 * usage: pass a object with field name as key,
 * "|" separated string as rules of the field
 * example: {"name":"required|minLength:5|maxLength:20"}
 * @param rules { [key: string]: string }
 */
export default function validator(rules: { [key: string]: string }) {
  let validatePassed = true;
  Object.keys(rules).forEach((key) => {
    const rulesOfAttr = rules[key].split("|");
    for (let rule of rulesOfAttr) {
      if (rule == "required") {
        if (!$(`[name|='${key}']`).val()) {
          $(`[name|='${key}']`).addClass("is-invalid");
          validatePassed = false;
          break;
        }
      } else if (rule.includes("maxLength")) {
        const value = $(`[name|='${key}']`).val() as string;
        if (
          value &&
          typeof value == "string" &&
          value.length > parseInt(rule.split(":")[1])
        ) {
          $(`[name|='${key}']`).addClass("is-invalid");
          validatePassed = false;
          break;
        }
      } else if (rule.includes("minLength")) {
        const value = $(`[name|='${key}']`).val();
        if (
          value &&
          typeof value == "string" &&
          value.length < parseInt(rule.split(":")[1])
        ) {
          $(`[name|='${key}']`).addClass("is-invalid");
          validatePassed = false;
          break;
        }
      } else if (rule.includes("regexp")) {
        const value = $(`[name|='${key}']`).val() as string;
        const regexp = new RegExp(rule.split(":")[1]);
        if (!regexp.test(value)) {
          $(`[name|='${key}']`).addClass("is-invalid");
          validatePassed = false;
          break;
        }
      } else if (rule.includes("max")) {
        const value = parseFloat($(`[name|='${key}']`).val() as string);
        if (
          value &&
          typeof value == "number" &&
          value > parseInt(rule.split(":")[1])
        ) {
          $(`[name|='${key}']`).addClass("is-invalid");
          validatePassed = false;
          break;
        }
      } else if (rule.includes("min")) {
        const value = parseFloat($(`[name|='${key}']`).val() as string);
        if (
          value &&
          typeof value == "number" &&
          value < parseInt(rule.split(":")[1])
        ) {
          $(`[name|='${key}']`).addClass("is-invalid");
          validatePassed = false;
          break;
        }
      } else if (rule.includes("same")) {
        const sourceValue = $(`[name|='${key}']`).val();
        const targetValue = $(`[name|='${rule.split(":")[1]}']`).val();
        if (sourceValue != targetValue) {
          $(`[name|='${key}']`).addClass("is-invalid");
          validatePassed = false;
          break;
        }
      } else if (rule == "email") {
        const value = $(`[name|='${key}']`).val() as string;
        if (value && typeof value == "string" && !isEmail(value)) {
          $(`[name|='${key}']`).addClass("is-invalid");
          validatePassed = false;
          break;
        }
      } else if (rule == "alpha") {
        const value = $(`[name|='${key}']`).val() as string;
        const alphaRegexp = new RegExp("^[a-zA-Z]*$");
        if (!alphaRegexp.test(value)) {
          $(`[name|='${key}']`).addClass("is-invalid");
          validatePassed = false;
          break;
        }
      }
    }
  });
  return validatePassed;
}
