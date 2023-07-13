import http from 'k6/http';
import { check, sleep } from 'k6';

// export let options = {
//     vus: 2, // Set the number of VUs to match the number of users
//     duration: '30s', // Duration of the test
//   };
  

export default function () {
  // Simulate concurrent requests for different operations
  const operations = ['create', 'update', 'delete'];

  // Define session IDs for users
  const sessionIDs = ['2cObVEq0fMMmXOHV0CGGiwZ2nWbJw4A7os37PaEg1wk=', 'CiKLAr1FGi7FNQiS34wbENjpc38nLsDMHqtroT1MsTo='];

  // Choose a random session ID from the array
//   const sessionID = sessionIDs[Math.floor(Math.random() * sessionIDs.length)];

  // Choose an operation randomly
  const operation = operations[0];
  sessionIDs.forEach((sessionID) => {
    // const seatNumber = 'C3'; // Use a specific seat number or generate dynamically
    // const payload = { show_id: 123, seat_number: seatNumber, seat_type_id: 1, version: 1 };


  // Simulate different seat numbers and show IDs
//   const seatNumber = `C-${Math.floor(Math.random() * 5) + 1}`; 
  const seatNumber="C-5"
  const showID = 1;

  switch (operation) {
    case 'create':
      createSeat(seatNumber, showID, sessionID);
      break;
    case 'update':
      updateSeat(seatNumber, showID, sessionID);
      break;
    case 'delete':
      deleteSeat(seatNumber, showID, sessionID);
      break;
    default:
      console.log('Invalid operation');
  }

})}

function createSeat(seatNumber, showID, sessionID) {
  const payload = { seat_number: seatNumber, show_id: showID };

  const res = http.post('http://localhost:9010/seats', JSON.stringify(payload), {
    headers: {
      'Content-Type': 'application/json',
      'Session-Id': sessionID,
    },
  });

  check(res, {
    'Seat created successfully': (r) => r.status === 201,
    'Failed to create seat': (r) => r.status === 500,
  });

  // Process the response or perform assertions if needed
  // ...
}

function updateSeat(seatNumber, showID, sessionID) {
  const queryParams = `seat_number=${seatNumber}&show_id=${showID}`;

  const res = http.put(`http://localhost:9010/seats?${queryParams}`, null, {
    headers: { 'Session-Id': sessionID },
  });

  check(res, {
    'Seat updated successfully': (r) => r.status === 200,
    'Seat not found': (r) => r.status === 404,
    'Failed to update seat': (r) => r.status === 500,
  });

  // Process the response or perform assertions if needed
  // ...
}

function deleteSeat(seatNumber, showID, sessionID) {
  const queryParams = `seat_number=${seatNumber}&show_id=${showID}`;

  const res = http.del(`http://localhost:9010/seats?${queryParams}`, null, {
    headers: { 'Session-Id': sessionID },
  });

  check(res, {
    'Seat deleted successfully': (r) => r.status === 200,
    'Seat not found': (r) => r.status === 404,
    'Failed to delete seat': (r) => r.status === 500,
  });

  // Process the response or perform assertions if needed
  // ...
}
