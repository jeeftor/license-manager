import React from 'react';

interface Props {
    name: string;
}

const HelloWorld: React.FC<Props> = ({ name }) => {
    return <div>Hello, {name}!</div>;
};

export default HelloWorld;
