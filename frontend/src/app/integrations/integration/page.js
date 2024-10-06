"use client";

import React, { useEffect, useState } from 'react';
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart";
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { BarChart, Bar, CartesianGrid, XAxis, PieChart, Pie, Label } from 'recharts';
import { Badge } from '@/components/ui/badge';
import { useRouter } from 'next/navigation';


import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

const integrationData = {
  name: "Temu FRA",
  createdAt: "2024-05-01",
  lastUpdate: "2024-10-05",
  totalTonnage: 2324,
  barChartData: [
    { month: 'May', shipments: 20000, tonnage: 100000 },
    { month: 'Aug', shipments: 25000, tonnage: 120000 },
    { month: 'Oct', shipments: 30000, tonnage: 150000 },
  ],
  pieChartData: [
    { type: "cleared", shipments: 27500, fill: "var(--color-cleared)" },
    { type: "blocked", shipments: 200, fill: "var(--color-blocked)" },
  ],
};

const chartConfig = {
  shipments: {
    label: "shipments",
  },
  cleared: {
    label: "Cleared Shipments",
    color: "#36b47E",
  },
  blocked: {
    label: "blocked Shipments",
    color: "#D35178",
  },
};

function PieChartComponent({ data }) {
  const totalshipments = React.useMemo(() => {
    return data.reduce((acc, curr) => acc + curr.shipments, 0);
  }, [data]);

  return (
    <div style={{ width: '100%', height: '100%', minWidth: '150px', minHeight: '150px' }}>
      <ChartContainer
        config={chartConfig}
        className="mx-auto aspect-square max-h-[250px]"
      >
        <PieChart>
          <ChartTooltip
            cursor={false}
            content={<ChartTooltipContent hideLabel />}
          />
          <Pie
            data={data}
            dataKey="shipments"
            nameKey="type"
            innerRadius={30}
            outerRadius={35}
            fill="#8884d8"
            paddingAngle={5}
            strokeWidth={5}
          >
            <Label
              content={({ viewBox }) => {
                if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                  return (
                    <text
                      x={viewBox.cx}
                      y={viewBox.cy}
                      textAnchor="middle"
                      dominantBaseline="middle"
                      fill="#fff"
                    >
                      {totalshipments}
                    </text>
                  );
                }
                return null;
              }}
            />
          </Pie>
        </PieChart>
      </ChartContainer>
    </div>
  );
}

export default function IntegrationPage() {
  const [isMounted, setIsMounted] = useState(false);
  const [tableRows, setTableRows] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);

  const specificRows = [
    {
      logisticsObject: "https://www.ecom-pipeline.net/12424",
      recipient: "Dr. Murat Yalcin",
      message: "Found HS-Code: 123456. Mismatch! Recommended HS-Code: 654321",
      badgeClass: "bg-yellow-500 hover:bg-yellow-500",
      badgeText: "HS Mismatch",
    },
    {
      logisticsObject: "https://www.ecom-pipeline.net/124224",
      recipient: "Benjamin Riege",
      message: "Found HS-Code: 234567. Possible Dangerous Goods!",
      badgeClass: "bg-red-500 hover:bg-red-500",
      badgeText: "Hidden DG",
    },
    {
      logisticsObject: "https://www.ecom-pipeline.net/122424",
      recipient: "Henk Mulder",
      message: "Found HS-Code: 345678. Mismatch! Recommended HS-Code: 876543",
      badgeClass: "bg-yellow-500 hover:bg-yellow-500",
      badgeText: "HS Mismatch",
    },
    {
      logisticsObject: "https://www.ecom-pipeline.net/124264",
      recipient: "Arnaud Lambert",
      message: "Reciever on Sanction List!",
      badgeClass: "bg-red-500 hover:bg-red-500",
      badgeText: "Sanction List",
    },
  ];

  useEffect(() => {
    setIsMounted(true);
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      setTableRows((prevRows) => {
        if (currentIndex < specificRows.length) {
          const newRow = specificRows[currentIndex];
          setCurrentIndex(currentIndex + 1);
          return [...prevRows, newRow];
        }
        clearInterval(interval);
        return prevRows;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [currentIndex]);

  if (!isMounted) {
    return null;
  }

  return (
    <div className="flex justify-center items-center min-h-screen">
      <div className="flex flex-col items-center gap-5">
        <h1 className="mt-16 mb-8 text-4xl">{integrationData.name}</h1>
        <div className="container p-5 rounded-lg inline-block rounded-lg">
          <div className="flex justify-between">
            <Card className="w-[450px] h-[370px] cursor-pointer hover:shadow-lg transition-shadow duration-300 rounded-lg">
              <CardHeader>
                <CardTitle className="text-2xl">Metrics</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="w-full">
                  <p><strong>active since:</strong> {integrationData.createdAt}</p>
                  <p><strong>last update:</strong> {integrationData.lastUpdate}</p>
                </div>
                <div className="flex justify-center gap-5 mt-5">
                  <div className="h-38 w-38">
                    <BarChart width={150} height={150} data={integrationData.barChartData}>
                      <CartesianGrid vertical={false} />
                      <XAxis dataKey="month" tickLine={false} tickMargin={10} axisLine={false} tickFormatter={(value) => value.slice(0, 3)} />
                      <Bar dataKey="shipments" fill="#A08FDB" radius={4} />
                    </BarChart>
                  </div>
                  <div className="h-38 w-38">
                    <PieChartComponent data={integrationData.pieChartData} />
                  </div>
                </div>
                <div className="mt-5 text-center">
                  <p><strong>Total Tonnage:</strong> {integrationData.totalTonnage}</p>
                </div>
              </CardContent>
            </Card>
            <Card className="w-[450px] h-[370px] cursor-pointer hover:shadow-lg transition-shadow duration-300 rounded-lg">
              <CardHeader>
                <CardTitle className="text-2xl">Stats</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="font-bold text-lg">total LO: <span className="text-green-500">10.5K</span></p>
                <p className="font-bold text-lg">total exceptions: <span className="text-red-400">27</span></p>
                <div className="flex flex-col">
                  <div className="relative">
                    <Badge
                        variant="secondary"
                        className="text-xl mt-8 max-w-xs cursor-pointer">
                          temu-fra@ecom-pipeline.net
                    </Badge>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
          <div className="flex justify-center gap-5 mt-5">
            <Card className="w-[1000px] h-[370px] cursor-pointer hover:shadow-lg transition-shadow duration-300 rounded-lg">
              <Table>
                <TableCaption>Data Enrichment and Exception Handling</TableCaption>
                <TableHeader>
                  <TableRow>
                    <TableHead>LogisticsObject</TableHead>
                    <TableHead>Recipient</TableHead>
                    <TableHead>Message</TableHead>
                    <TableHead className="text-right"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {tableRows.map((row, index) => (
                    <TableRow key={index}>
                      <TableCell>{row.logisticsObject}</TableCell>
                      <TableCell>{row.recipient}</TableCell>
                      <TableCell>{row.message}</TableCell>
                      <TableCell className="text-right"><Badge className={row.badgeClass}>{row.badgeText}</Badge></TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}