/**
 * Submits query to the server
 * @param config - Frontend configuration
 * @returns {function({preparedQuery?: *}): Promise<unknown>}
 */
const createSubmitToApi = function(config) {
    return function submitQueryToAPI({preparedQuery}) {
        return new Promise((resolve, reject) => {
            $.post(`${config.apiHost}:${config.apiPort}/${config.routes.queryRoute}`,{
                query: preparedQuery
            }).fail((err) => {
                reject(err.statusText);
            }).done((res) => {
                resolve(res);
            })
        });
    }
};