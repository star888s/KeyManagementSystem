import React from 'react';
import { MdDone } from 'react-icons/md';
import { AiOutlineClose } from 'react-icons/ai';

interface ScheduleModalProps {
  showModal: boolean;
  setShowModal: (showModal: boolean) => void;
  newSchedule: any; // Replace with your actual type
  handleIdChange: (e: React.ChangeEvent<HTMLSelectElement>) => void;
  handleTimeChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleInputChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  handleSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  infoList: any[]; // Replace with your actual type
}

const ScheduleModal: React.FC<ScheduleModalProps> = ({
  showModal,
  setShowModal,
  newSchedule,
  handleIdChange,
  handleTimeChange,
  handleInputChange,
  handleSubmit: originalHandleSubmit,
  infoList,
}) => {
  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const now = new Date();
    const startTime = new Date(newSchedule.startTime);
    const endTime = new Date(newSchedule.endTime);

    if (startTime < now || endTime < now) {
      alert('現在時刻より後の値を指定してください。');
      return;
    }

    if (startTime >= endTime) {
      alert('開始時刻は終了時刻より前の値を指定してください。');
      return;
    }

    originalHandleSubmit(e);
  };
  if (!showModal) {
    return null;
  }

  return (
    <div
      style={{
        position: 'fixed',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        backgroundColor: '#EFF1E6',
        padding: '20px',
        border: '1px solid #333',
        borderRadius: '10px',
        zIndex: 9999,
      }}
    >
      <form onSubmit={handleSubmit}>
        <p>部屋名</p>
        <select
          name='id'
          value={newSchedule.id}
          onChange={handleIdChange}
          style={{
            boxShadow: '0px 0px 10px rgba(0, 0, 0, 0.5)',
            borderRadius: '5px',
          }}
        >
          {infoList.map((info) => (
            <option value={info.id}>{info.name}</option>
          ))}
        </select>
        <p>開始時刻</p>
        <input
          type='datetime-local'
          name='startTime'
          onChange={handleTimeChange}
          style={{
            boxShadow: '0px 0px 10px rgba(0, 0, 0, 0.5)',
            borderRadius: '5px',
          }}
        />
        <br />
        <p>終了時刻</p>
        <input
          type='datetime-local'
          name='endTime'
          onChange={handleTimeChange}
          style={{
            boxShadow: '0px 0px 10px rgba(0, 0, 0, 0.5)',
            borderRadius: '5px',
          }}
        />
        <br />
        <p>メモ欄</p>
        <input
          type='text'
          name='memo'
          maxLength={20}
          onChange={handleInputChange}
          placeholder='メモを記載...'
          style={{
            boxShadow: '0px 0px 10px rgba(0, 0, 0, 0.5)',
            borderRadius: '5px',
          }}
        />
        <br></br>
        <div style={{ textAlign: 'right' }}>
          <button type='submit'>
            <MdDone size={30} />
          </button>
        </div>
      </form>
      <button
        onClick={() => setShowModal(false)}
        style={{
          position: 'absolute',
          right: '10px',
          top: '10px',
          background: 'transparent',
          border: 'none',
          fontSize: '1.5em',
        }}
      >
        <AiOutlineClose size={30} />
      </button>
    </div>
  );
};

export default ScheduleModal;
