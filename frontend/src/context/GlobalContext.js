import { createContext, useState } from 'react';

export const GlobalContext = createContext();

export const GlobalProvider = ({ children }) => {
  const [schema, setSchema] = useState(null);
  const [username, setUsername] = useState(null);
  const [anotherState, setAnotherState] = useState(null); // Add more states as needed

  return (
    <GlobalContext.Provider value={{ schema, setSchema, anotherState, setAnotherState, username, setUsername }}>
      {children}
    </GlobalContext.Provider>
  );
};