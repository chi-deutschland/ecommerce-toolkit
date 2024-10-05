
"use client";
import React, { useContext } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { GlobalContext } from '@/context/GlobalContext'; // Adjust the import path as needed


const cardData = [
  { title: "Confirm Schema", content: "Schema has been confirmed." },
  { title: "Card 22", content: "Content for card 2." },
  { title: "Card 3", content: "Content for card 3." },
  { title: "Card 4", content: "Content for card 4." },
  { title: "Card 5", content: "Content for card 5." },
  { title: "Card 6", content: "Content for card 6." },
  { title: "Card 7", content: "Content for card 7." },
  { title: "Card 8", content: "Content for card 8." },
  { title: "Card 9", content: "Content for card 9." },
  { title: "Card 10", content: "Content for card 10." }
];

const bgColors = [
  "bg-gray-200",
  "bg-blue-200",
  "bg-green-200",
  "bg-yellow-200",
  "bg-red-200",
  "bg-purple-200",
  "bg-pink-200",
  "bg-indigo-200",
  "bg-teal-200",
  "bg-orange-200"
];

const ConfirmSchemaPage = () => {
  const { schema } = useContext(GlobalContext);
  console.log(schema);
  return (
    <div className="flex justify-center items-center min-h-screen">
      <div className="flex flex-col items-center gap-5">
        <h1 className="mt-16 mb-8 text-4xl">Confirm the Schema</h1>
        <div className="container p-5 inline-block rounded-lg">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
            {cardData.map((card, index) => (
              <Card key={index} className="w-[200px] h-[300px]">
                <CardHeader className={`rounded-t-lg ${bgColors[index % bgColors.length]} text-gray-800`}>
                  <CardTitle>{card.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p>{card.content}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ConfirmSchemaPage;