'use client';
import './globals.css';
import { useSession } from 'next-auth/react';
import { useRouter } from 'next/navigation';
import React, { useEffect } from 'react';
import { Login } from './components/AuthButton';

export default function Home() {
  const { data: session } = useSession();
  const router = useRouter();

  useEffect(() => {
    if (session) {
      router.push('/dashbord');
    }
  });
  return (
    <main className='flex h-screen justify-center items-center'>
      <div className='text-center'>
        <div>
          <h1 className='text-4xl font-bold my-10 mx-auto'>KeyManagementSystem</h1>
        </div>
        <Login></Login>
      </div>
    </main>
  );
}
