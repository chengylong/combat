// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
    uint256 public count;

    event CountedTo(uint256 newCount);

    function increment() public {
        count += 1;
        emit CountedTo(count);
    }

    function getCount() public view returns (uint256) {
        return count;
    }
}