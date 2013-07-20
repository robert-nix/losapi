# losapi

Todo: pretty this up

    base - /api/ (nginx proxy)
    
    all: limit, offset (limit cap: 500)
    
    messages: all sort received_-1 : /messages?
    by user=?                                  user={user}
    by channel=?                               channel={channel}
    by command=?                               command={command}
    received range:                            start={time}
    end requires start to have effect          end={time}
    then filter regex/i                        match={string}
    
    statuses: all sort timestamp_-1 :
    by channel=?                     /channel/:channel?
    then by timestamp range                            start={time}
     requires start                                    end={time}
    
    capped collection info: /cutoff
    returns dates for the oldest documents in messages and statuses
    {
      "messages": ...,
      "statuses": ...
    }
