'use server';
import { Schedule } from '../model'; // Replace with your actual type

const api_key = process.env.API_KEY || '';
const host = process.env.URL || '';

const deleteSchedule = async (schedule: Schedule) => {
  const headers = { 'X-Api-Key': api_key };
  const url = host + '/delete_schedule' || '';

  const response = await fetch(url, {
    method: 'POST',
    headers: headers,
    mode: 'cors',
    credentials: 'same-origin',
    body: JSON.stringify([schedule]),
  });

  if (!response.ok) {
    console.log(`HTTP error! status: ${response.status}`);
    return false;
  }

  return true;
};

export default deleteSchedule;
