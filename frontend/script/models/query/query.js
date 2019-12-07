/**
 * Validates and makes query object
 * @param validate - Validate function
 * @returns {function(*=): Readonly<{preparedQuery: *}>}
 */
makeQuery = function(validate) {
    return function(queryString) {
        let { error } = validate(queryString);
        if (error) {
            throw new Error(error);
        }
        return Object.freeze({
            preparedQuery: queryString
        });
    }
}