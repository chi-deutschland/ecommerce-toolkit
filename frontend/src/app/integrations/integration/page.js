import React from 'react';

const tableData = [
  { col1: "Data 1", col2: "Data 2", col3: "Data 3", col4: "Data 4", col5: "Data 5", col6: "Data 6", col7: "Data 7", col8: "Data 8", col9: "Data 9", col10: "Data 10", col11: "Data 11", col12: "Data 12", col13: "Data 13" },
  { col1: "Data 14", col2: "Data 15", col3: "Data 16", col4: "Data 17", col5: "Data 18", col6: "Data 19", col7: "Data 20", col8: "Data 21", col9: "Data 22", col10: "Data 23", col11: "Data 24", col12: "Data 25", col13: "Data 26" },
  // Add more rows as needed
];

const ConfirmSchemaPage = () => {
    return (
      <div>
        <h1 className="mt-16 mb-8 text-4xl text-center">Welcome James Bond</h1>
        <div className="flex justify-center">
          <div className="container mx-auto p-5 rounded-lg">
            <div className="overflow-x-auto">
              <table className="min-w-full bg-white rounded-lg">
                <thead>
                  <tr>
                    {Array.from({ length: 13 }, (_, i) => (
                      <th key={i} className="px-4 py-2 border-b-2 border-gray-300 text-left text-sm font-semibold text-gray-700">
                        Column {i + 1}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody className="rounded-lg">
                  {tableData.map((row, rowIndex) => (
                    <tr
                      key={rowIndex}
                      className={`transition-colors duration-200 ${rowIndex % 2 === 0 ? 'bg-gray-200' : 'bg-gray-300'} hover:bg-gray-400`}
                    >
                      {Object.values(row).map((cell, cellIndex) => (
                        <td
                          key={cellIndex}
                          className={`px-4 py-2 border-b border-gray-300 text-sm text-gray-700 ${
                            rowIndex === tableData.length - 1 && cellIndex === 0 ? 'rounded-bl-lg' : ''
                          } ${
                            rowIndex === tableData.length - 1 && cellIndex === Object.values(row).length - 1 ? 'rounded-br-lg' : ''
                          }`}
                        >
                          {cell}
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>

      </div>
    );
  };
  
  export default ConfirmSchemaPage;