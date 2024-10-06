
"use client";
import React, { useState, useContext } from 'react';
import { Card, CardHeader, CardTitle, CardContent, Button } from '@/components/ui/card';
import { GlobalContext } from '@/context/GlobalContext';
import { useRouter } from 'next/navigation';
import Select from 'react-select';
import axios from 'axios';


const ConfirmSchemaPage = () => {
  const { schema, setSchema } = useContext(GlobalContext);
  const [selectedValues, setSelectedValues] = useState({});
  const [loading, setLoading] = useState(false);
  const router = useRouter();



  const handleSubmit = async () => {
    setLoading(true);

    try {
      const response = await axios.post('http://35.246.224.221/pipeline', schema);
      console.log('Response:', response.data);
    } catch (error) {
      console.error('Error:', error);
    } finally {
      setLoading(false);
    }
    setLoading(true);
    setTimeout(() => {
      router.push('/integrations/one_record');
    }, 2500);
  };

  const dropDownItems = Object.values(schema)
    .filter(item => item.content)
    .map((item, index) => ({
      value: item.content,
      label: item.content,
    }));

  const customStyles = {
    control: (provided) => ({
      ...provided,
      backgroundColor: '#333',
      color: '#fff',
    }),
    singleValue: (provided) => ({
      ...provided,
      color: '#fff',
    }),
    menu: (provided) => ({
      ...provided,
      backgroundColor: '#333',
    }),
    option: (provided, state) => ({
      ...provided,
      backgroundColor: state.isSelected ? '#555' : '#333',
      color: '#fff',
      '&:hover': {
        backgroundColor: '#444',
      },
    }),
  };

  const handleSelectChange = (selectedOption, cardKey) => {
    setSelectedValues(prevState => ({
      ...prevState,
      [cardKey]: selectedOption.value,
    }));

    setSchema(prevSchema => ({
      ...prevSchema,
      [cardKey]: {
        ...prevSchema[cardKey],
        content: selectedOption.value,
      },
    }));
  };

  console.log('schema:', schema);

  return (
    <div className="flex justify-center min-h-screen">
      <div className="flex flex-col items-center gap-5">
        <div className="flex items-center justify-center mt-16 mb-8">
          <div className="flex items-center">
            <div className="flex items-center justify-center w-10 h-10 rounded-full bg-green-600 text-xl">✓</div>
              <div className="w-32 h-1 bg-green-600 mx-2"></div>
              <div className="flex items-center justify-center w-10 h-10 rounded-full bg-green-600 text-xl">2</div>
              <div className="w-32 h-1 bg-gray-300 mx-2"></div>
              <div className="flex items-center justify-center w-10 h-10 rounded-full bg-gray-300 text-xl">3</div>
            </div>
          </div>
        <div className="container p-5 inline-block rounded-lg relative">
          {loading && (
            <div className="absolute inset-0 flex justify-center items-end bg-gray-500 bg-opacity-50 z-50">
              <span className="loader mb-80"></span>
            </div>
          )}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
            {Object.entries(schema).map(([key, card]) => (
              key !== 'headers' && (
                <Card key={key} className="w-[300px] h-[350px]">
                  <CardHeader className={`h-[75px] rounded-t-lg bg-gray-400 bg-opacity-55 text-gray-200`}>
                    <CardTitle>{card.title}</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="flex flex-col justify-between h-full">
                      <p className="font-bold text-lg mt-4">mapped value:</p>
                      <p  dangerouslySetInnerHTML={{ __html: card.content }}></p>
                      <p className="font-bold text-lg mt-8">corrected value:</p>
                      <Select
                        options={dropDownItems}
                        styles={customStyles}
                        onChange={(selectedOption) => handleSelectChange(selectedOption, key)}
                      />
                    </div>
                  </CardContent>
                </Card>
              )
            ))}
          </div>
          <div className="flex justify-end">
            <button
              onClick={handleSubmit}
              className="h-[50px] mt-8 px-4 py-2 bg-blue-500 text-white rounded"
              disabled={loading}
            >
              {loading ? 'Submitting...' : 'Submit Schema'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ConfirmSchemaPage;