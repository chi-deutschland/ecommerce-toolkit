"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card } from '@/components/ui/card';
import { useAuth } from '@/context/AuthContext';

export default function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const { setIsLoggedIn } = useAuth();
  const router = useRouter();

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log('Username:', username);
    console.log('Password:', password);
    setIsLoggedIn(true);
    router.push('/integrations');
  };

  return (
    <div>
      <div className="flex justify-center items-center mt-8 scale-up mb-16">
        <div className="mt-1 block-left ml-4 w-2 h-2 bg-gray-500"></div>
        <div className="mb-2 block-left ml-4 w-2 h-2 bg-gray-600"></div>
        <div className="mt-3 block-left ml-4 w-2 h-2 bg-gray-700"></div>
        <div className="block-left ml-4 w-2 h-2 bg-gray-700"></div>
        <div className="ml-4 w-2 h-9 bg-gray-500 pipe"></div>
        <div className="w-32 h-8 bg-gray-200 pipe text-dark text-gray-900">E-Comm Pipeline</div>
        <div className="mr-4 w-2 h-9 bg-gray-500 pipe"></div>
        <div className="uniform-block"></div>
      </div>
      <div className="flex justify-center">
        <Card className="w-80 p-5">
          <div className="mb-5">
            <h2 className="text-2xl font-semibold">Login</h2>
          </div>
          <form onSubmit={handleSubmit} className="flex flex-col">
            <Label className="mb-2">
              Username:
              <Input 
                type="text" 
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="mt-1 block w-full"
              />
            </Label>
            <Label className="mb-2">
              Password:
              <Input 
                type="password" 
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="mt-1 block w-full"
              />
            </Label>
            <Button type="submit" className="mt-4 bg-blue-500 text-white py-2 px-4 rounded">
              Login
            </Button>
          </form>
        </Card>
      </div>
    </div>
  );
}