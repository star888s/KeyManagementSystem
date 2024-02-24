import { format } from 'date-fns';
import { AiOutlineClose, AiOutlineDelete } from 'react-icons/ai';
import { Schedule } from '../model';

import React from 'react';

interface ScheduleDetailModalProps {
  selectedSchedule: any;
  setSelectedSchedule: (selectedSchedule: any) => void;
  deleteSchedule: (selectedSchedule: any) => void;
}

const ScheduleDetailModal: React.FC<ScheduleDetailModalProps> = ({
  selectedSchedule,
  setSelectedSchedule,
  deleteSchedule,
}) => {
  if (!selectedSchedule) {
    return null;
  }

  // ISO形式の日付をJavaScriptのDateオブジェクトに変換
  const startTime = new Date(selectedSchedule.startTime);
  const endTime = new Date(selectedSchedule.endTime);

  // DateオブジェクトをHH:mm形式にフォーマット
  const formattedStartTime = format(startTime, 'HH:mm');
  const formattedEndTime = format(endTime, 'HH:mm');

  return (
    <div
      className='fixed top-0 left-0 w-screen h-screen bg-black bg-opacity-50'
      onClick={() => setSelectedSchedule(null)}
    >
      <div
        className='card absolute p-5 border border-gray-700'
        style={{
          top: '40%',
          left: '35%',
          transform: 'translate(-50%, -50%)',
          width: '20%',
          height: '30%',
        }}
        onClick={(e) => e.stopPropagation()}
      >
        <button className='absolute top-0 right-0 p-2' onClick={() => setSelectedSchedule(null)}>
          <AiOutlineClose size={30} />
        </button>
        <h2 className='text-2xl'>{selectedSchedule.name}</h2>
        <div className='time'>
          {formattedStartTime} - {formattedEndTime}
        </div>
        <h3 style={{ fontFamily: 'Courier New, monospace' }}>メモ</h3>
        <textarea
          style={{
            fontFamily: 'inherit',
            maxWidth: '100%', // 親要素の幅に合わせる
            maxHeight: '100%', // 親要素の高さに合わせる
            height: '100px',
            whiteSpace: 'pre-wrap',
            wordWrap: 'break-word',
            border: '1px solid #000',
            borderRadius: '5px',
            boxShadow: '2px 2px 5px rgba(0, 0, 0, 0.2)',
            overflow: 'auto',
          }}
          readOnly
        >
          {selectedSchedule.memo}
        </textarea>
        <br></br>
        <br></br>
        <button className='absolute bottom-2 right-2 p-2' onClick={() => deleteSchedule(selectedSchedule)}>
          <AiOutlineDelete size={30} />
        </button>
      </div>
    </div>
  );
};

export default ScheduleDetailModal;
