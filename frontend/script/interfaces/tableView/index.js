/**
 * Draws table dynamically based on returned results
 * @param rootElementId - Element to be appended upon
 * @param dataArray - Returned results
 * @returns {HTMLParagraphElement}
 */
const drawTable = function(rootElementId, dataArray) {
    if (dataArray.length === 0) {
        let noElementsParagraph = document.createElement('p');
        noElementsParagraph.innerHTML = 'Empty result returned';
        return noElementsParagraph;
    }

    const columns = Object.keys(dataArray[0]);

    let tableElement = document.createElement('table');
    let tableHeadElement = document.createElement('thead');
    let tableBodyElement = document.createElement('tbody');

    tableElement.style.border = '1px solid black';

    for (const columnName of columns) {
        let headerColumn = document.createElement('th');
        headerColumn.style.border = '1px solid black';
        headerColumn.innerHTML = columnName;
        tableHeadElement.appendChild(headerColumn);
    }
    tableElement.appendChild(tableHeadElement);

    for (const row of dataArray) {
        let tableRow = document.createElement('tr');
        for (const columnName of columns) {
            let rowData = document.createElement('td');
            rowData.style.border = '1px solid black';
            rowData.innerHTML = row[columnName];
            tableRow.appendChild(rowData);
        }
        tableBodyElement.appendChild(tableRow);
    }
    tableElement.appendChild(tableBodyElement);

    const root = document.getElementById(rootElementId);
    root.innerHTML = '';
    root.appendChild(tableElement);
}