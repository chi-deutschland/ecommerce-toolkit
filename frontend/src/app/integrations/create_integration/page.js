"use client";

import { useContext, useState} from 'react';
import { useDropzone } from 'react-dropzone';
import { Input } from "@/components/ui/input";
import { GlobalContext } from '@/context/GlobalContext'; // Adjust the import path as needed
import { useRouter } from 'next/navigation';

const CreateIntegrationPage = () => {
  const { schema, setSchema } = useContext(GlobalContext);
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const [file, setFile] = useState(null);

  const onDrop = (acceptedFiles) => {
    console.log(acceptedFiles);
    if (acceptedFiles.length > 0) {
      setFile(acceptedFiles[0]);
    }
  };

  const handleUpload = async () => {
    if (file) {
        setIsLoading(true);
        const formData = new FormData();
        formData.append('file', file);

        try {
            const response = await fetch('http://35.246.224.221/schema', {
                method: 'POST',
                body: formData,
            });

            const data = await response.json();
            setSchema(data);
            setTimeout(() => {
              router.push('/integrations/confirm_schema');
            }, 1000);
        } catch (error) {
            console.error('Error:', error);
        }
    } else {
        console.error('No file selected');
    }
};

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    accept: '.xls',
    onDrop,
  });

  return (
    <div className="flex flex-col">
      <div className="flex items-center justify-center mt-16">
        <div className="flex items-center">
          <div className="flex items-center justify-center w-10 h-10 rounded-full bg-green-600 text-xl">1</div>
          <div className="w-32 h-1 bg-gray-300 mx-2"></div>
          <div className="flex items-center justify-center w-10 h-10 rounded-full bg-gray-300 text-xl">2</div>
          <div className="w-32 h-1 bg-gray-300 mx-2"></div>
          <div className="flex items-center justify-center w-10 h-10 rounded-full bg-gray-300 text-xl">3</div>
        </div>
      </div>
      <div className="flex justify-center items-center min-h-screen position-relative">
        <div className="flex flex-col items-center gap-5">
          <div className="bg-gray-100 w-[900px] p-5 custom-dropzone-container rounded-lg">
            <div className="flex flex-col items-center gap-4 custom-dropzone-container">
              <div className="w-full max-w-md">
                <Input id="name" name="name" type="text" placeholder="Name of Integration" className="mt-1 block w-full bg-white text-gray-800" />
              </div>
              <div
                {...getRootProps()}
                className={`flex justify-center items-center h-96 bg-gray-200 w-full border-2 border-dashed ${
                  isDragActive ? 'border-blue-500' : 'border-gray-300 custom-dropzone'
                }`}
              >
                <input {...getInputProps()} />
                {isDragActive ? (
                  <p>Drop the files here ...</p>
                ) : (
                  <p>Place your file here (only .xls)</p>
                )}
              </div>
              <button
                onClick={handleUpload}
                className="mt-4 px-4 py-2 bg-blue-500 text-white rounded"
              >
                Upload File
              </button>
            </div>
          </div>
        </div>
        {isLoading && (
          <div className="absolute inset-0 flex justify-center items-center bg-gray-500 bg-opacity-50 z-50">
            <span className="loader"></span>
          </div>
        )}
      </div>
    </div>
  );
};

export default CreateIntegrationPage;