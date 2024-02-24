'use server';
import { Schedule } from '../model';

const api_key = process.env.API_KEY || '';
const host = process.env.URL || '';

const fetchData = async () => {
  const headers = { 'X-Api-Key': api_key };
  const url = host + '/get_schedule' || '';

  const response = await fetch(url, {
    method: 'GET',
    headers: headers,
    mode: 'cors',
    credentials: 'same-origin',
    cache: 'no-store',
  });
  const data = await response.json();

  // dataのサイズを確認
  console.log(data.Items.length);

  // Sort schedules by startTime in ascending order
  const sortedSchedules = data.Items.sort((a: Schedule, b: Schedule) => {
    return new Date(a.startTime).getTime() - new Date(b.startTime).getTime();
  });

  console.log(sortedSchedules.length);

  return sortedSchedules;
};

export default fetchData;
