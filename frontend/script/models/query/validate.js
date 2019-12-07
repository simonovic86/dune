/**
 * Validates basic SQL query with regular expression
 * @param queryString - Query string to be validated
 * @returns {{error: *}}
 */
function validateQuery(queryString) {
    let error = null;
    let re = /\s*select\s*([*]|(\s*\w+(\s*[as|]\s*\w+)?([,]\s*\w+\s*)*))\s*(from\s*.*\s*(where\s*.*)?)?[;]?/gmi;

    if (!queryString.match(re)) {
        error = "Invalid SQL query";
    }
    return { error };
}