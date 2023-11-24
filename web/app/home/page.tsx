'use client';
import { useSession } from 'next-auth/react';
import { Logout } from '../components/AuthButton';
import { useRouter } from 'next/navigation';
import React, { useEffect } from 'react';

export default function Home() {
  const { data: session } = useSession();
  const router = useRouter();
  useEffect(() => {
    if (!session) {
      router.push('/');
    }
  });
  return (
    <main>
      <div className='max-w-[85rem] px-4 py-10 sm:px-6 lg:px-8 lg:py-14 mx-auto'>
        <div className='flex flex-col'>
          <div className='-m-1.5 overflow-x-auto'>
            <div className='p-1.5 min-w-full inline-block align-middle'>
              <div className='bg-white border border-gray-200 rounded-xl shadow-sm overflow-hidden dark:bg-slate-900 dark:border-gray-700'>
                <div className='px-6 py-4 grid gap-3 md:flex md:justify-between md:items-center border-b border-gray-200 dark:border-gray-700'>
                  <div>
                    <h2 className='text-xl font-semibold text-gray-800 dark:text-gray-200'></h2>
                    Hello!
                    <Logout></Logout>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
