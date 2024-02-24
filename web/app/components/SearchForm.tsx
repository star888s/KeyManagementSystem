import React from 'react';
import { MdCancel } from 'react-icons/md';

interface SearchFormProps {
  searchTerm: string;
  setSearchTerm: (searchTerm: string) => void;
}

const SearchForm: React.FC<SearchFormProps> = ({ searchTerm, setSearchTerm }) => {
  return (
    <div className='flex items-center justify-center'>
      <input
        className='rounded py-2 px-4'
        type='text'
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        placeholder='Search...'
      />
      <button className='ml-4' onClick={() => setSearchTerm('')}>
        <MdCancel size={30} />
      </button>
    </div>
  );
};

export default SearchForm;
