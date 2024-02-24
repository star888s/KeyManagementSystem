import React from 'react';
import { FiPlus } from 'react-icons/fi';

interface AddScheduleButtonProps {
  setShowModal: (showModal: boolean) => void;
}

const AddScheduleButton: React.FC<AddScheduleButtonProps> = ({ setShowModal }) => {
  return (
    <button onClick={() => setShowModal(true)}>
      <FiPlus size={30} />
    </button>
  );
};

export default AddScheduleButton;
