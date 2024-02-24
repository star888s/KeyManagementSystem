'use server';

const api_key = process.env.API_KEY || '';
const host = process.env.URL || '';

const fetchInfo = async () => {
  const headers = { 'X-Api-Key': api_key };
  const url = host + '/get_info' || '';

  const response = await fetch(url, {
    method: 'GET',
    headers: headers,
    mode: 'cors',
    credentials: 'same-origin',
  });
  const data = await response.json();

  const infoItems = data.Items;

  return infoItems;
};

export default fetchInfo;
