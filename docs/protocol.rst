.. |--| unicode:: U+2013   .. en dash
.. |---| unicode:: U+2014  .. em dash, trimming surrounding whitespace
   :trim:

Entangle protocol
=================

This document describes revision 1 of the Entangle protocol. Note that revision to the Entangle protocol will only be performed with major version changes of Entangle, effectively meaning a protocol freeze per major version.


Data interchange format
-----------------------

The Entangle protocol uses `MessagePack <http://msgpack.org/>`_ as the data interchange format. Please note that Entangle provides data types not native to MessagePack which are therefore to be interpreted during data deserialization.


Wire protocol
-------------

The wire protocol of Entangle is effectively transport agnostic. However, implementations are allowed to expect ordering as a property of the underlying transport as the wire protocol does not provide a transport layer multiplexing facility.

All messages sent over the wire in either direction are effectively MessagePack encoded arrays. Please note that as the wire protocol does not provide transport multiplexing, messages must be sent in full. The first two elements are generic to all messages:

::

   [<opcode>, <message ID>, ..]

``opcode``
   **Opcode** |--| *uint8*

   Opcode for the specific message. Can be one of the following, indicating the type of message:

   +-----------+---------------------------------+
   | Value     | Message type                    |
   +===========+=================================+
   | ``0x00``  | `Request`_                      |
   +-----------+---------------------------------+
   | ``0x01``  | `Notification`_                 |
   +-----------+---------------------------------+
   | ``0x02``  | `Response`_                     |
   +-----------+---------------------------------+
   | ``0x03``  | `Exception`_                    |
   +-----------+---------------------------------+
   | ``0x04``  | `Notification acknowledgement`_ |
   +-----------+---------------------------------+
   | ``0x7f``  | `Compressed message`_           |
   +-----------+---------------------------------+

``message ID``
   **Message ID** |--| *uint32*

   Message ID unique to the connection. Used by clients and servers to identify requests and responses. As far less than :math:`2^{32}-1` outstanding requests per connection are expected at one time in reality, it is considered safe to wrap the value around to :math:`0` on overflow.


Messages
--------

Request
~~~~~~~

A request is a message requesting the execution of a remote method. A request must under normal working conditions be responded to with either a `response`_ or an `exception`_.

The structure of a request is as follows:

::

   [<opcode>, <message ID>, <method>, <arguments>, <request trace>]

``method``
   **Remote method name** |--| *string*

   Name of the remote method to be executed.

``arguments``
   **Method arguments** |--| *array*

   Method arguments as an array of zero or more arbitrary values. Cannot be ``nil``.

``request trace``
   **Request execution trace** |--| *bool*

   Whether an execution trace is wanted.


Notification
~~~~~~~~~~~~

A notification as a message requesting the execution of a remote method without receiving for the result. Note however, that to properly detect and handle connection and communication issues, a message must be sent in the response to a notification. The message can be either a `notification acknowledgement`_ or an `exception`_.

The structure of a notification is as follows:

::

   [<opcode>, <message ID>, <method>, <arguments>]

``method``
   **Remote method name** |--| *string*

   Name of the remote method to be executed.

``arguments``
   **Method arguments** |--| *array*

   Method arguments as an array of zero or more arbitrary values. Cannot be ``nil``.


Response
~~~~~~~~

A request is a message requesting the execution of a remote method. The structure of a request is as follows:

::

   [<opcode>, <message ID>, <result>, <trace>]

``result``
   **Result of method** |--| *arbitrary*

   Arbitrary result value. Can be ``nil``.

``trace``
   **Execution trace** |--| *trace*

   Execution trace. Can be ``nil`` indicating that no trace was requested or
   generated. Note that it is considered bad practise for a server to not
   return a trace upon request.


Exception
~~~~~~~~~

A request is a message requesting the execution of a remote method. The structure of a request is as follows:

::

   [<opcode>, <message ID>, <definition>, <name>, <description>, <request trace>]

``definition``
   **Exception source definition** |--| *string*

   Name of the definition in which the exception is declared.

``name``
   **Exception name** |--| *string*

   Name of the exception.

``description``
   **Excption description** |--| *string*

   Description of the exception that occured.

``trace``
   **Execution trace** |--| *trace*

   Execution trace. Can be ``nil`` indicating that no trace was requested or
   generated. Note that it is considered bad practise for a server to not
   return a trace upon request. However, expect the trace to be ``nil`` if the
   exception occurs prior to executing the requested method.


Notification acknowledgement
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

A notification acknowledgement is a message indicating that a notification has been successfully received. The struture of a notification acknowledgement is as follows:

::

   [<opcode>, <message ID>]


Compressed message
~~~~~~~~~~~~~~~~~~

A compressed message is essentially a wrapper for a message. Thus, the compressed data contained in the message is in itself a message. It is a requirement that both the compressed message and the contained message have the same message ID. The structure of a compressed message is as follows:

::

   [<opcode>, <message ID>, <compression method>, <compressed data>]

``compression method``
   **Compression method** |--| *uint8*

   Compression method used to compress the data. Can be one of the following:

   +----------+-----------------------------------------------+
   | Value    | Compression method                            |
   +==========+===============================================+
   | ``0x00`` | `Snappy <https://code.google.com/p/snappy/>`_ |
   +----------+-----------------------------------------------+

``compressed data``
   **Compressed data** |--| *binary*

   Data compressed using the indicated compression method. The data itself contains a message.


Error handling
--------------

The Entangle protocol opts for an agressive failure approach in order to prevent communication issues resulting in unexpected behaviour in production environments. Thus, Entangle categorises errors into two: `recoverable errors`_ and `unrecoverable errors`_.


Recoverable errors
~~~~~~~~~~~~~~~~~~




Unrecoverable errors
~~~~~~~~~~~~~~~~~~~~


